import React, { FunctionComponent } from "react";

const Button: FunctionComponent<{
    label: string;
}> = function ({ label }) {
    return <button aria-label={label}>{label}</button>
};
export default Button;
