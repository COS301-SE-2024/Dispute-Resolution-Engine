import renderer from 'react-test-renderer';
import { Card } from '../card';

it('renders correctly', () => {
  const tree = renderer
    .create(<Card />)
    .toJSON();
  expect(tree).toMatchSnapshot();
});