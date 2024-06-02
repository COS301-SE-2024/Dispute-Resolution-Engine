import renderer from 'react-test-renderer';
import { ScrollArea } from '../scroll-area';

it('renders correctly', () => {
  const tree = renderer
    .create(<ScrollArea />)
    .toJSON();
  expect(tree).toMatchSnapshot();
});
