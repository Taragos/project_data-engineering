import { env } from "$env/dynamic/private";

export async function load() {
  const res = await fetch(
    `http://${env.DATA_ACCESS_SERVICE_HOST}:${env.DATA_ACCESS_SERVICE_PORT}/api/all/`
  );
  const results = await res.json();

  return { results, s3Url: `http://${env.S3_EXTERNAL_ENDPOINT}/${env.S3_BUCKET}/` };
}
