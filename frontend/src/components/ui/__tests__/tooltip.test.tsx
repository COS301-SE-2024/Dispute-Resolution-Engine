import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import '@testing-library/jest-dom';
import { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider } from '../tooltip';
import { Button } from '../button';

describe('Tooltip', () => {
  it('renders correctly', () => {
    render(
      <TooltipProvider data-testid="tt">
        <Tooltip >
          <TooltipTrigger>
            Hover me
          </TooltipTrigger>
          <TooltipContent >Test tooltip</TooltipContent>
        </Tooltip>
      </TooltipProvider>
    );

    expect(screen.queryByText('Test tooltip')).not.toBeInTheDocument();
    userEvent.hover(screen.getByText('Hover me'));
    expect(screen.getByText('Hover me')).toBeInTheDocument();
    
  });
});