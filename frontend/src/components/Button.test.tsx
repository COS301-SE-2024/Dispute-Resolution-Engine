import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import Button from "./Button";

describe("Button", () => {
    it("renders a button", () => {
        render(<Button label="label" />);

        const heading = screen.getByRole('button');
        expect(heading).toBeInTheDocument();
    })
    it("has an accessibility label", () => {
        const label = "My label";
        render(<Button label={label} />);

        const elem = screen.getByLabelText(label);
        expect(elem).toBeInTheDocument();
    })
    it("shows label visibly", () => {
        const label = "My label";
        render(<Button label={label} />);

        const elem = screen.getByText(label);
        expect(elem).toBeInTheDocument();
    })
})
