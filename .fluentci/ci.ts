import {
  fmt,
  test,
  build,
} from "https://pkg.fluentci.io/go_pipeline@v0.6.1/mod.ts";

await fmt();
await test();
await build();
