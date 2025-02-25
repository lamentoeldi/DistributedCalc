"use client"

import styles from "../../public/expressions.module.css";
import {useEffect, useState} from "react";

export interface expression {
    id: string,
    result: number,
    status: string
}

const Expression = (exp: expression) => {
    const [color, setColor] = useState("")

    useEffect(()=> {
        switch (exp.status.toLowerCase()) {
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
        <tr className={styles.Exp}>
            <td className={styles.Td}>
                <span>{exp.id}</span>
            </td>
            <td className={styles.Td}>
                <span>{exp.result}</span>
            </td>
            <td className={styles.Td}>
                <span className={color}>{exp.status}</span>
            </td>
        </tr>
    )
}

export default Expression