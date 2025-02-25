"use client"

import Navbar from "@/components/Navbar";
import styles from "../../../public/main.module.css"
import Calculate from "@/components/Calculate";

export default function Page() {
    return (
        <div className={styles.Wrapper}>
            <div className={styles.Main}>
                <Navbar/>
                <div className={styles.Frame}>
                    <h1>Start new calculation</h1>
                    <Calculate/>
                </div>
            </div>
        </div>
    );
}
