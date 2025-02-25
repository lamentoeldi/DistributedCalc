import Navbar from "@/components/Navbar";
import styles from "../../../public/main.module.css"
import Find from "@/components/Find";

export default function Page() {
    return (
        <div className={styles.Wrapper}>
            <div className={styles.Main}>
                <Navbar/>
                <div className={styles.Frame}>
                    <h1>Find expression by ID</h1>
                    <Find/>
                </div>
            </div>
        </div>
    );
}
