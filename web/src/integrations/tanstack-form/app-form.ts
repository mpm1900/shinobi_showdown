import { createFormHook, createFormHookContexts, } from "@tanstack/react-form";

const { fieldContext, formContext, useFieldContext, useFormContext } =
  createFormHookContexts()

const { useAppForm, withForm } = createFormHook({
  fieldComponents: {},
  formComponents: {},
  fieldContext,
  formContext
})

export { fieldContext, formContext, useFieldContext, useFormContext, useAppForm, withForm }
