# IKOU

A React SSR/SSG framework built in Go. This project is mainly an experiment with Go, React and server side rendering and I don't recommend you actually
use this in a production environment.

## Features

- File based routing
- Server side rendering
- Static site generation
- API routes
- Middleware
- TailwindCSS support

## Getting started

### Installation

### React File Structure

#### Pages

### Backend File Structure

The two forms of backend routes are API routes and Entry routes.
These are both defined in the `routes` directory and into their respective `api` and `entry` subdirectories.
Like the React pages, the routes are also defined by the file structure.

#### API Routes

API routes allow you to create RESTful endpoints in your application.
This defaults to the `/api` path but can be changed in the configuration.

#### Entry Routes

Entry routes allow you to run Go code on the server before rendering the page.
The entry route handler function also let's you return data to be passed as props to the page.

### Middleware

### Configuration

## Roadmap

## Contributing

## License

MIT

See the [LICENSE](LICENSE) file for more information.
