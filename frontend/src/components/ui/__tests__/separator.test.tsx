import renderer from 'react-test-renderer';
import { Separator } from '../separator';

it('renders correctly', () => {
  const tree = renderer
    .create(<Separator />)
    .toJSON();
  expect(tree).toMatchSnapshot();
});
