import Client, { connect } from "https://sdk.fluentci.io/v0.1.9/mod.ts";
import {
  fmt,
  test,
  build,
} from "https://pkg.fluentci.io/go_pipeline@v0.5.1/mod.ts";

function pipeline(src = ".") {
  connect(async (client: Client) => {
    await fmt(client, src);
    await test(client, src);
    await build(client, src);
  });
}

pipeline();
