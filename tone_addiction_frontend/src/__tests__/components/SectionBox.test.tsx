import { render, screen } from '@testing-library/react';
import { ChakraProvider } from "@chakra-ui/react";
import {ReactElement} from "react";
import {SectionBox} from "../../components/SectionBox.tsx";

const renderWithChakra = (ui: ReactElement) => {
  return render(<ChakraProvider>{ui}</ChakraProvider>);
};

describe('SectionBox', () => {
  it('renders with the correct title', () => {
    const title = 'Test Title';
    renderWithChakra(<SectionBox title={title} />);

    const headingElement = screen.getByRole('heading', { name: title });
    expect(headingElement).toBeInTheDocument();
  });

  it('renders the Divider component', () => {
    renderWithChakra(<SectionBox title="Another Title" />);

    const dividerElement = screen.getByRole('separator');
    expect(dividerElement).toBeInTheDocument();
  });

  it('renders the AbsoluteCenter component', () => {
    renderWithChakra(<SectionBox title="Centered Title" />);

    const absoluteCenterElement = screen.getByText('Centered Title').parentElement;
    expect(absoluteCenterElement).toHaveStyle('position: absolute');
  });
});
