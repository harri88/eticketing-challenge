import { render, screen } from '@testing-library/react';
import App from './App';

test('renders Event E-Ticketing heading', () => {
  render(<App />);
  const headingElement = screen.getByText(/Event E-Ticketing/i);
  expect(headingElement).toBeInTheDocument();
});
