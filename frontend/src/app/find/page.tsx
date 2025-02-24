import Navbar from "@/components/Navbar";
import styles from "../../../public/main.module.css"

export default function Page() {
    return (
        <div className={styles.Wrapper}>
            <div className={styles.Main}>
                <Navbar/>
                <div className={styles.Frame}>
                    <h2>Find expression by ID</h2>
                </div>
            </div>
        </div>
    );
}
