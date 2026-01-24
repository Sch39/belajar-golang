import {
  BookOpenIcon,
  CodeBracketIcon,
  RocketLaunchIcon,
} from "@heroicons/react/24/outline";

export default function LandingPage() {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-900">
      {/* HERO */}
      <section className="mx-auto max-w-6xl px-6 py-24 text-center">
        <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
          My Kasir REST API Gw
        </h1>
        <p className="mt-4 text-lg text-gray-600">
          Web API hasil belajar Golang — clean architecture, RESTful,
          dan siap dipakai untuk eksperimen frontend.
        </p>

        <div className="mt-8 flex justify-center gap-4">
          <a
            href="/swagger"
            className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-6 py-3 text-white shadow hover:bg-blue-700"
          >
            <BookOpenIcon className="h-5 w-5" />
            View Documentation
          </a>

          <a
            href="https://github.com/Sch39/belajar-golang"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-6 py-3 text-gray-700 hover:bg-gray-100"
          >
            <CodeBracketIcon className="h-5 w-5" />
            Source Code
          </a>
        </div>
      </section>

      {/* FEATURES */}
      <section className="mx-auto max-w-6xl px-6 py-16">
        <div className="grid gap-8 sm:grid-cols-3">
          <Feature
            icon={<RocketLaunchIcon className="h-6 w-6 text-blue-600" />}
            title="Fast & Simple"
            description="Built with net/http and clean service-repository pattern."
          />
          <Feature
            icon={<CodeBracketIcon className="h-6 w-6 text-blue-600" />}
            title="RESTful JSON API"
            description="Consistent response format, clear error handling, and validation."
          />
          <Feature
            icon={<BookOpenIcon className="h-6 w-6 text-blue-600" />}
            title="Swagger Documentation"
            description="Auto-generated OpenAPI docs for easy exploration."
          />
        </div>
      </section>

      {/* SAMPLE */}
      <section className="mx-auto max-w-4xl px-6 py-16">
        <h2 className="text-2xl font-semibold">Sample Request</h2>

        <pre className="mt-4 overflow-x-auto rounded-lg bg-gray-900 p-4 text-sm text-gray-100">
{`curl -X GET https://api.example.com/api/categories

{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Kopi",
      "description": "Menu kopi"
    }
  ]
}`}
        </pre>
      </section>

      {/* FOOTER */}
      <footer className="border-t bg-white">
        <div className="mx-auto max-w-6xl px-6 py-8 text-sm text-gray-500">
          <div className="flex flex-col gap-2 sm:flex-row sm:justify-between">
            <span>© {new Date().getFullYear()} Golang Learning API</span>
            <div className="flex gap-4">
              <a href="/swagger" className="hover:text-gray-900">
                Swagger
              </a>
              <a
                href="https://github.com/Sch39/belajar-golang"
                target="_blank"
                rel="noopener noreferrer"
                className="hover:text-gray-900"
              >
                GitHub
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}

function Feature({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="rounded-xl bg-white p-6 shadow-sm">
      <div className="mb-4">{icon}</div>
      <h3 className="text-lg font-semibold">{title}</h3>
      <p className="mt-2 text-sm text-gray-600">{description}</p>
    </div>
  );
}
