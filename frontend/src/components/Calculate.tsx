import styles from "../../public/find.module.css";
import React, {useState} from "react";
import style from "../../public/expression.module.css"

const Calculate = () => {
    const [id, setId] = useState<number>(0)
    const [err, setErr] = useState("")
    const [value, setValue] = useState("");

    const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setValue(event.target.value);
    };

    const start = async (exp: string) => {
        const res = await fetch("/api/v1/calculate", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({expression: exp})
        })

        if (!res.ok) {
            setErr(res.statusText)
            return
        }

        const data = await res.json()

        setId(data.id)
    }

    return (
        <div>
            <input
                placeholder={"Enter expression"}
                value={value}
                onChange={handleChange}
                className={styles.Input}
            />
            <button
                onClick={() => start(value)}
                className={styles.Send}
            >
                Calculate
            </button>
            {
                id != 0 ? (
                    <div className={style.ID}>
                        ID: {id}
                    </div>
                ) : err && (
                    <h3>{err}</h3>
                )
            }
        </div>
    )
}

export default Calculate