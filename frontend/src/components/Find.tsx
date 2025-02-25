"use client"

import {expression} from "@/components/ExpressionTable";
import React, {useEffect, useState} from "react";
import styles from "../../public/find.module.css"
import Expression from "@/components/Expression";

const Find = () => {
    const [exp, setExp] = useState<expression>()
    const [err, setErr] = useState("")
    const [value, setValue] = useState("");
    const [isLoading, setIsLoading] = useState(false)

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setValue(event.target.value);
    };

    const find = async (id: string) => {
        if (id === "") {
            setErr("Expression ID must not be empty");
            return
        }

        setIsLoading(true)
        const r = await fetch("api/v1/expressions/" + id);
        setIsLoading(false)
        if (!r.ok) {
            setErr(r.statusText);
            return
        }
        const data = await r.json();
        setExp(data.expression);
    }

    return (
        <div>
            <input
                placeholder={"Enter expression ID"}
                value={value}
                onChange={handleChange}
                className={styles.Input}
            />
            <button
                onClick={() => find(value)}
                className={styles.Send}
            >Find</button>
            {
                isLoading ? (
                        <h1>Loading...</h1>
                ) : exp ? (
                    <Expression id={exp.id} result={exp.result} status={exp.status}/>
                ) : err && (
                    <h3>{err}</h3>
                )
            }
        </div>
    )
}

export default Find