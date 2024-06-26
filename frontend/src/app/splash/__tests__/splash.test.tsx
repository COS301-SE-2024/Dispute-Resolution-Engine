import { render, screen, fireEvent } from '@testing-library/react';

import "@testing-library/jest-dom";
import SplashPage from "@/app/splash/splash";


describe('Splash Page', () => {
  it('renders without crashing', () => {
    render(<SplashPage />);
    expect(screen.getByText('Here to streamline your ADR process')).toBeInTheDocument();
  });
});