import Navbar from "@/components/Navbar";
import styles from "../../public/main.module.css"
import Expression from "@/components/Expressions";
import Expressions from "@/components/Expressions";

export default function Page() {
  return (
      <div className={styles.Wrapper}>
          <div className={styles.Main}>
            <Navbar/>
            <div className={styles.Frame}>
                <h1>Expressions</h1>
                <Expressions/>
            </div>
          </div>
      </div>
  );
}
