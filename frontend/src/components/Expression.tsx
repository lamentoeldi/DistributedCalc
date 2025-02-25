"use client"

import {expression} from "@/components/ExpressionTable";
import styles from "../../public/expression.module.css"
import {useEffect, useState} from "react";

const Expression = (e: expression) => {
    const [color, setColor] = useState("")

    useEffect(()=> {
        switch (e.status.toLowerCase()) {
            case "completed":
                setColor(styles.Green)
                break
            case "failed":
                setColor(styles.Red)
                break
            case "pending":
                setColor(styles.Yellow)
        }
    }, [])

    return (
        <div className={styles.Wrapper}>
            <span>ID: {e.id}</span>
            <span>Result: {e.result}</span>
            <span>Status:
                <span className={color}>{e.status.charAt(0).toUpperCase() + e.status.slice(1)}</span>
            </span>
        </div>
    )
}

export default Expression