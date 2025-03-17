

// NewCmdEdit creates the `edit` command
func NewCmdEdit(f cmdutil.Factory, ioStreams genericiooptions.IOStreams) *cobra.Command {
	o := editor.NewEditOptions(editor.NormalEditMode, ioStreams)
	cmd := &cobra.Command{
		Use:                   "edit (RESOURCE/NAME | -f FILENAME)",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Edit a resource on the server"),
		Long:                  editLong,
		Example:               editExample,
		ValidArgsFunction:     completion.ResourceTypeAndNameCompletionFunc(f),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, args, cmd))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())  # 이 부분이 핵심 로직 시작
		},
	}

	// bind flag structs
	o.RecordFlags.AddFlags(cmd)
	o.PrintFlags.AddFlags(cmd)

	usage := "to use to edit the resource"
	cmdutil.AddFilenameOptionFlags(cmd, &o.FilenameOptions, usage)
	cmdutil.AddValidateFlags(cmd)
	cmd.Flags().BoolVarP(&o.OutputPatch, "output-patch", "", o.OutputPatch, "Output the patch if the resource is edited.")
	cmd.Flags().BoolVar(&o.WindowsLineEndings, "windows-line-endings", o.WindowsLineEndings,
		"Defaults to the line ending native to your platform.")
	cmdutil.AddFieldManagerFlagVar(cmd, &o.FieldManager, "kubectl-edit")
	cmdutil.AddApplyAnnotationVarFlags(cmd, &o.ApplyAnnotation)
	cmdutil.AddSubresourceFlags(cmd, &o.Subresource, "If specified, edit will operate on the subresource of the requested object.")
	return cmd
}




// Run performs the execution
func (o *EditOptions) Run() error {
	edit := NewDefaultEditor(editorEnvs())
	// editFn is invoked for each edit session (once with a list for normal edit, once for each individual resource in a edit-on-create invocation)
	editFn := func(infos []*resource.Info) error {
		var (
			results  = editResults{}
			original = []byte{}
			edited   = []byte{}
			file     string
			err      error
		)

		containsError := false
		// loop until we succeed or cancel editing
		for {
			// get the object we're going to serialize as input to the editor
			var originalObj runtime.Object
			switch len(infos) {
			case 1:
				originalObj = infos[0].Object
			default:
				l := &unstructured.UnstructuredList{
					Object: map[string]interface{}{
						"kind":       "List",
						"apiVersion": "v1",
						"metadata":   map[string]interface{}{},
					},
				}
				for _, info := range infos {
					l.Items = append(l.Items, *info.Object.(*unstructured.Unstructured))
				}
				originalObj = l
			}

			// generate the file to edit
			buf := &bytes.Buffer{}
			var w io.Writer = buf
			if o.WindowsLineEndings {
				w = crlf.NewCRLFWriter(w)
			}

			if o.editPrinterOptions.addHeader {
				results.header.writeTo(w, o.EditMode)
			}

			if !containsError {
				if err := o.extractManagedFields(originalObj); err != nil {
					return preservedFile(err, results.file, o.ErrOut)
				}

				if err := o.editPrinterOptions.PrintObj(originalObj, w); err != nil {
					return preservedFile(err, results.file, o.ErrOut)
				}
				original = buf.Bytes()
			} else {
				// In case of an error, preserve the edited file.
				// Remove the comments (header) from it since we already
				// have included the latest header in the buffer above.
				buf.Write(cmdutil.ManualStrip(edited))
			}

			// launch the editor
			editedDiff := edited


			✔️ kubectl edit은 내부적으로 LaunchTempFile을 호출하여 /tmp/kubectl-edit-XXXX.yaml 같은 임시 파일을 생성함.
			✔️ 공식 문서에는 직접적으로 명시되어 있지는 않지만, 소스 코드에서 확인 가능.
			✔️ kubectl edit 실행 후 ls /tmp/kubectl-edit-* 하면 실제 파일이 생성되는 것을 볼 수 있음. 🚀

			edited, file, err = edit.LaunchTempFile(fmt.Sprintf("%s-edit-", filepath.Base(os.Args[0])), o.editPrinterOptions.ext, buf)
			if err != nil {
				return preservedFile(err, results.file, o.ErrOut)
			}

			// If we're retrying the loop because of an error, and no change was made in the file, short-circuit
			if containsError && bytes.Equal(cmdutil.StripComments(editedDiff), cmdutil.StripComments(edited)) {
				return preservedFile(fmt.Errorf("%s", "Edit cancelled, no valid changes were saved."), file, o.ErrOut)
			}
			// cleanup any file from the previous pass
			if len(results.file) > 0 {
				os.Remove(results.file)
			}
			klog.V(4).Infof("User edited:\n%s", string(edited))

			// Apply validation
			schema, err := o.f.Validator(o.ValidationDirective)
			if err != nil {
				return preservedFile(err, file, o.ErrOut)
			}
			err = schema.ValidateBytes(cmdutil.StripComments(edited))
			if err != nil {
				results = editResults{
					file: file,
				}
				containsError = true
				fmt.Fprintln(o.ErrOut, results.addError(apierrors.NewInvalid(corev1.SchemeGroupVersion.WithKind("").GroupKind(),
					"", field.ErrorList{field.Invalid(nil, "The edited file failed validation", fmt.Sprintf("%v", err))}), infos[0]))
				continue
			}

			// Compare content without comments
			if bytes.Equal(cmdutil.StripComments(original), cmdutil.StripComments(edited)) {
				os.Remove(file)
				fmt.Fprintln(o.ErrOut, "Edit cancelled, no changes made.")
				return nil
			}

			lines, err := hasLines(bytes.NewBuffer(edited))
			if err != nil {
				return preservedFile(err, file, o.ErrOut)
			}
			if !lines {
				os.Remove(file)
				fmt.Fprintln(o.ErrOut, "Edit cancelled, saved file was empty.")
				return nil
			}

			results = editResults{
				file: file,
			}

			// parse the edited file
			updatedInfos, err := o.updatedResultGetter(edited).Infos()
			if err != nil {
				// syntax error
				containsError = true
				results.header.reasons = append(results.header.reasons, editReason{head: fmt.Sprintf("The edited file had a syntax error: %v", err)})
				continue
			}

			// not a syntax error as it turns out...
			containsError = false
			updatedVisitor := resource.InfoListVisitor(updatedInfos)

			// we need to add back managedFields to both updated and original object
			if err := o.restoreManagedFields(updatedInfos); err != nil {
				return preservedFile(err, file, o.ErrOut)
			}
			if err := o.restoreManagedFields(infos); err != nil {
				return preservedFile(err, file, o.ErrOut)
			}

			// need to make sure the original namespace wasn't changed while editing
			if err := updatedVisitor.Visit(resource.RequireNamespace(o.CmdNamespace)); err != nil {
				return preservedFile(err, file, o.ErrOut)
			}

			// iterate through all items to apply annotations
			if err := o.visitAnnotation(updatedVisitor); err != nil {
				return preservedFile(err, file, o.ErrOut)
			}

			switch o.EditMode {
			case NormalEditMode:
				err = o.visitToPatch(infos, updatedVisitor, &results)
			case ApplyEditMode:
				err = o.visitToApplyEditPatch(infos, updatedVisitor)
			case EditBeforeCreateMode:
				err = o.visitToCreate(updatedVisitor)
			default:
				err = fmt.Errorf("unsupported edit mode %q", o.EditMode)
			}
			if err != nil {
				return preservedFile(err, results.file, o.ErrOut)
			}

			// Handle all possible errors
			//
			// 1. retryable: propose kubectl replace -f
			// 2. notfound: indicate the location of the saved configuration of the deleted resource
			// 3. invalid: retry those on the spot by looping ie. reloading the editor
			if results.retryable > 0 {
				fmt.Fprintf(o.ErrOut, "You can run `%s replace -f %s` to try this update again.\n", filepath.Base(os.Args[0]), file)
				return cmdutil.ErrExit
			}
			if results.notfound > 0 {
				fmt.Fprintf(o.ErrOut, "The edits you made on deleted resources have been saved to %q\n", file)
				return cmdutil.ErrExit
			}

			if len(results.edit) == 0 {
				if results.notfound == 0 {
					os.Remove(file)
				} else {
					fmt.Fprintf(o.Out, "The edits you made on deleted resources have been saved to %q\n", file)
				}
				return nil
			}

			if len(results.header.reasons) > 0 {
				containsError = true
			}
		}
	}

