import SplashPage from "./splash";
import SplashHeader from "@/app/splash/splash-header";

export default function Splash() {
  return (


    <div className={"relative"}>
      <div className="fixed w-full z-40"><SplashHeader /></div>
      <main className="pt-20">
        <SplashPage />
      </main>
    </div>

  );
}
