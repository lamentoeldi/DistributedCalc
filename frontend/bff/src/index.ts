import {Elysia, error, t } from "elysia";

const backend_host = process.env.BACKEND_HOST ?? "localhost"
const backend_port = process.env.BACKEND_PORT ?? "8080"

const port = process.env.PORT ?? "3000"

const logger = new Elysia()
    .onRequest(({ request }) => {
        console.log(`received request ${request.url}`)
    })
    .onError(({ code, error }) => {
        console.log("err: ", code, error)
    })

const auth = new Elysia()
    .post(
        "/bff/api/v1/register",
        async ({ body }) => {
            const res = await fetch(`http://${backend_host}:${backend_port}/api/v1/register`, {
                method: "POST",
                body: JSON.stringify(body)
            })

            if (!res.ok) {
                error(res.status)
            }
        },
        {
            body: t.Object({
                login: t.String(),
                password: t.String()
            })
        }
    )
    .post(
        "/bff/api/v1/login",
        async ({ body, cookie: { access_token, refresh_token } }) => {
            if (!refresh_token) {
                error(401)
                return
            }

            const res = await fetch(`http://${backend_host}:${backend_port}/api/v1/login`, {
                method: "POST",
                body: JSON.stringify(body)
            })

            if (!res.ok) {
                error(res.status)
                return
            }

            const data = await res.json()

            access_token.value = data.access_token
            access_token.httpOnly = true

            refresh_token.value = data.refresh_token
            refresh_token.httpOnly = true
        },
        {
            body: t.Object({
                login: t.String(),
                password: t.String()
            })
        }
    )
    .get(
        "/bff/api/v1/authorize",
        async ({ cookie: { access_token, refresh_token } } ) => {
            const res = await fetch(`http://${backend_host}:${backend_port}/api/v1/authorize`, {
                method: "GET",
                headers: {
                    'Authorization': `Bearer ${access_token.value}`,
                    'Refresh-Token': `${refresh_token.value}`,
                    'Accept': 'application/json',
                }
            })

            if (!res.ok) {
                error(res.status)
                return
            }

            console.log("log res data: ", await res.json())
            const data: {
                user_id: string,
                username: string
            } = await res.json()

            return data
        }
    )

const calculator = new Elysia()
    .post(
        "/bff/api/v1/calculate",
        async ({ body, cookie: { access_token, refresh_token } }) => {
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

            const data: {
                id: string
            } = await res.json()

            return data
        },
        {
            body: t.Object({
                expression: t.String()
            })
        }
    )
    .get(
        "/bff/api/v1/expressions",
        async ({ query: { cursor, limit }, cookie: { access_token, refresh_token } }) => {
            let url = `http://${backend_host}:${backend_port}/api/v1/expressions`

            if (cursor && limit) {
                url += `?cursor=${cursor}&limit=${limit}`
            } else if (cursor) {
                url += `?cursor=${cursor}`
            } else if (limit) {
                url += `?limit=${limit}`
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

            const data: {
                expressions: {
                   id: string
                   status: string
                   result: number
                }[]
            } = await res.json()

            return data
        },
        {
            query: t.Object({
                cursor: t.String(),
                limit: t.Number()
            })
        }
    )
    .get("/bff/api/v1/expressions/:id", async ({ params: { id }, cookie: { access_token, refresh_token }}) => {
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

        const data: {
            expressions: {
                id: string
                status: string
                result: number
            }
        } = await res.json()

        return data
    })

const app = new Elysia()
    .use(logger)
    .use(calculator)
    .use(auth)
    // staticPlugin was supposed to be used here, though it does not support SPAs
    // implemented simple workaround
    // check https://github.com/elysiajs/elysia-static/issues/22#issue-2010404588 for more info
    .get('/*', async ({ path }) => {
        const staticFile = Bun.file(`./.dist/${path}`);
        const fallBackFile = Bun.file('./.dist/index.html');
        return (await staticFile.exists()) ? staticFile : fallBackFile;
    })
    .listen(port);

export type App = typeof app

console.log(
  `Elysia is running at http://${app.server?.hostname}:${app.server?.port}`
);