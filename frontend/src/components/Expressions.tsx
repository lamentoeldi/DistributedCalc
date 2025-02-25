"use client"

import styles from "../../public/expressions.module.css"
import Expression from "@/components/ExpressionTable";
import {useEffect, useState} from "react";
import {expression} from "@/components/ExpressionTable";

const Expressions = () => {
    const [exp, setExp] = useState<expression[]>()

    useEffect(() => {
        (async () => {
           const res = await fetch("api/v1/expressions")
            const data = await res.json();
            setExp(data.expressions)
        })()
    }, [])

    if (!exp || exp.length < 1) {
        return <h3>No expressions are being processed</h3>
    }

    return (
        <table className={styles.ExpTb}>
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Result</th>
                    <th>Status</th>
                </tr>
            </thead>
            <tbody>
            {   exp && (
                    exp.map((e: expression)=> (
                        <Expression
                            key={e.id}
                            id={e.id}
                            result={e.result}
                            status={e.status.charAt(0).toUpperCase() + e.status.slice(1)}
                        />
                    ))
                )
            }
            </tbody>
        </table>
    )
}

export default Expressions