// TextField.test.tsx
import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { useForm, FormProvider } from 'react-hook-form';
import TextField from '../text-field';
import { SignupData } from '@/app/signup/signup-form';

const Wrapper: React.FC = ({ children }) => {
  const methods = useForm<SignupData>();
  return <FormProvider {...methods}>{children}</FormProvider>;
};

describe('TextField component', () => {
  it('renders with correct label and placeholder', () => {
    render(
      <Wrapper>
        <TextField<SignupData> name="username" label="Username" type="text" />
      </Wrapper>
    );

    // Check if the label is rendered
    expect(screen.getByLabelText('Username')).toBeInTheDocument();
    
    // Check if the input with placeholder is rendered
    const input = screen.getByPlaceholderText('Username');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('type', 'text');
  });

  it('renders password input when type is password', () => {
    render(
      <Wrapper>
        <TextField<SignupData> name="password" label="Password" type="password" />
      </Wrapper>
    );

    // Check if the input with placeholder is rendered
    const input = screen.getByPlaceholderText('Password');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('type', 'password');
  });
});
