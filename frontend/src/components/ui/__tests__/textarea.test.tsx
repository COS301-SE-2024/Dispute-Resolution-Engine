import renderer from 'react-test-renderer';
import { Textarea } from '../textarea';

it('renders correctly', () => {
  const tree = renderer
    .create(<Textarea />)
    .toJSON();
  expect(tree).toMatchSnapshot();
});
