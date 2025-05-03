import {Elysia, error} from "elysia";
import { staticPlugin } from "@elysiajs/static"

const backend_host = process.env.BACKEND_HOST ?? "localhost"
const backend_port = process.env.BACKEND_PORT ?? "8080"

const port = process.env.PORT ?? "3000"

const bff = new Elysia()
    .post("/api/v1/register", async ({ body }) => {
        const res = await fetch(`http://${backend_host}:${backend_port}/api/v1/register`, {
            method: "POST",
            body: JSON.stringify(body)
        })

        if (!res.ok) {
            error(res.status)
        }
    })
    .post("/api/v1/login", async ({ body, cookie: { access_token, refresh_token } }) => {
        const res = await fetch(`http://${backend_host}:${backend_port}/api/v1/login`, {
            method: "POST",
            body: JSON.stringify(body)
        })

        if (!res.ok) {
            error(res.status)
            return
        }

        const rb = await res.json()

        access_token.value = rb.access_token
        access_token.httpOnly = true

        refresh_token.value = rb.refresh_token
        refresh_token.httpOnly = true
    })
    .post("/api/v1/calculate", async ({ body, cookie: { access_token, refresh_token } }) => {
        const res = await fetch(`http://${backend_host}:${backend_port}/api/v1/calculate`, {
            method: "POST",
            body: JSON.stringify(body),
            headers: {
                'Authorization': `Bearer ${access_token.value}`,
                'Refresh-Token': `${refresh_token.value}`,
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })

        if (!res.ok) {
            error(res.status)
            return
        }

        const access = res.headers.get("Access-Token")
        if (access !== null) {
            access_token.value = access
            access_token.httpOnly = true
        }

        const refresh = res.headers.get("Refresh-Token")
        if (refresh !== null) {
            refresh_token.httpOnly = true
            refresh_token.httpOnly = true
        }

        const data = await res.json()

        return { data }
    })
    .get("/api/v1/expressions", async ({ query: { cursor, limit }, cookie: { access_token, refresh_token } }) => {
        let url = `http://${backend_host}:${backend_port}/api/v1/expressions`

        if (cursor && limit) {
            url += `?cursor=${cursor}&limit=${limit}`
        } else if (cursor) {
            url += `?cursor=${cursor}`
        } else if (limit) {
            url += `?limit=${limit.length}`
        }

        const res = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${access_token.value}`,
                'Refresh-Token': `${refresh_token.value}`,
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })

        if (!res.ok) {
            error(res.status)
            return
        }

        const access = res.headers.get("Access-Token")
        if (access !== null) {
            access_token.value = access
            access_token.httpOnly = true
        }

        const refresh = res.headers.get("Refresh-Token")
        if (refresh !== null) {
            refresh_token.httpOnly = true
            refresh_token.httpOnly = true
        }

        const data = await res.json()

        return { data }
    })
    .get("/api/v1/expressions/:id", async ({ params: { id }, cookie: { access_token, refresh_token }}) => {
        const res = await fetch(`http://${backend_host}:${backend_port}/api/v1/expressions/${id}`, {
            headers: {
                'Authorization': `Bearer ${access_token.value}`,
                'Refresh-Token': `${refresh_token.value}`,
                'Content-Type': 'application/json',
                'Accept': 'application/json',
            }
        })

        if (!res.ok) {
            error(res.status)
            return
        }

        const access = res.headers.get("Access-Token")
        if (access !== null) {
            access_token.value = access
            access_token.httpOnly = true
        }

        const refresh = res.headers.get("Refresh-Token")
        if (refresh !== null) {
            refresh_token.httpOnly = true
            refresh_token.httpOnly = true
        }

        const data = await res.json()

        return { data }
    })

const app = new Elysia()
    .use(staticPlugin({
        prefix: '/',
        alwaysStatic: true
    }))
    .use(bff)
    .listen(port);

console.log(
  `Elysia is running at http://${app.server?.hostname}:${app.server?.port}`
);