import renderer from 'react-test-renderer';
import { Label } from '../label';

it('renders correctly', () => {
  const tree = renderer
    .create(<Label />)
    .toJSON();
  expect(tree).toMatchSnapshot();
});
