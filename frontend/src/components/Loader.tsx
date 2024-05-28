import Image from 'next/image';
import Ralph from '../../public/ralph.png';

export default function Loader() {
  return (
    <div className="flex justify-center items-center h-screen">
	  <div className="animate-spin bg-white rounded-full h-32 w-32 border-8 border-black overflow-hidden">
		<Image src={Ralph} alt="Ralph" className="h-32 w-32 mt-6" />
	  </div>
	</div>
  );
}
