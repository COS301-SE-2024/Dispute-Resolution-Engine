import * as React from "react";
import { useFormState } from "react-dom";

const FormContext = React.createContext<any | undefined>(undefined);

export function createFormContext<T>(
  formName: string,
  action: (state: Awaited<T | undefined>, data: FormData) => T | Promise<T>
): [
  React.Context<T | undefined>,
  React.ForwardRefExoticComponent<
    React.HTMLAttributes<HTMLFormElement> & React.RefAttributes<HTMLFormElement>
  >
] {
  const form = React.forwardRef<HTMLFormElement, React.HTMLAttributes<HTMLFormElement>>(
    (props, ref) => {
      const [state, formAction] = useFormState(action, undefined);

      return (
        <FormContext.Provider value={state}>
          <form action={formAction} {...props} ref={ref} />
        </FormContext.Provider>
      );
    }
  );
  form.displayName = formName;

  return [FormContext, form];
}
