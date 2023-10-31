import Client, { connect } from "../../deps.ts";

export enum Job {
  test = "test",
  fmt = "fmt",
  build = "build",
}

export const exclude = ["vendor", ".git"];

export const test = async (src = ".") => {
  await connect(async (client: Client) => {
    const context = client.host().directory(src);
    const ctr = client
      .pipeline(Job.test)
      .container()
      .from("golang:latest")
      .withDirectory("/app", context, { exclude })
      .withWorkdir("/app")
      .withMountedCache("/go/pkg/mod", client.cacheVolume("go-mod"))
      .withMountedCache("/root/.cache/go-build", client.cacheVolume("go-build"))
      .withExec(["go", "test", "-v", "./..."]);
    const result = await ctr.stdout();

    console.log(result);
  });
  return "Done";
};

export const fmt = async (src = ".") => {
  await connect(async (client: Client) => {
    const context = client.host().directory(src);
    const ctr = client
      .pipeline(Job.fmt)
      .container()
      .from("golang:latest")
      .withDirectory("/app", context, { exclude })
      .withMountedCache("/go/pkg/mod", client.cacheVolume("go-mod"))
      .withMountedCache("/root/.cache/go-build", client.cacheVolume("go-build"))
      .withWorkdir("/app")

      .withExec(["go", "fmt", "./..."]);
    const result = await ctr.stdout();

    console.log(result);
  });
  return "Done";
};

export const build = async (src = ".") => {
  await connect(async (client: Client) => {
    const context = client.host().directory(src);
    const ctr = client
      .pipeline(Job.build)
      .container()
      .from("golang:latest")
      .withExec(["apt-get", "update"])
      .withExec([
        "apt-get",
        "install",
        "-y",
        "libsqlite3-mod-spatialite",
        "build-essential",
      ])
      .withMountedCache("/assets", client.cacheVolume("gh-release-assets"))
      .withDirectory("/app", context, { exclude })
      .withWorkdir("/app")
      .withMountedCache("/go/pkg/mod", client.cacheVolume("go-mod"))
      .withMountedCache("/root/.cache/go-build", client.cacheVolume("go-build"))
      .withExec(["go", "build", "-o", "bin/mada"])
      .withEnvVariable("TAG", Deno.env.get("TAG") || "latest")
      .withExec([
        "sh",
        "-c",
        "tar cvf /assets/mada_${TAG}_x86_64-unknown-linux-gnu.tar.gz bin/mada",
      ])
      .withExec([
        "sh",
        "-c",
        "shasum -a 256 /assets/mada_${TAG}_x86_64-unknown-linux-gnu.tar.gz > /assets/mada_${TAG}_x86_64-unknown-linux-gnu.tar.gz.sha256",
      ]);
    await ctr.stdout();
  });
  return "Done";
};

export type JobExec = (src?: string) => Promise<string>;

export const runnableJobs: Record<Job, JobExec> = {
  [Job.test]: test,
  [Job.fmt]: fmt,
  [Job.build]: build,
};

export const jobDescriptions: Record<Job, string> = {
  [Job.test]: "Run tests",
  [Job.fmt]: "Format code",
  [Job.build]: "Build binary",
};
