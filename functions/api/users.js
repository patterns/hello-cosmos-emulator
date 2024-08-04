/**
 * GET /api/users
 */
export async function onRequestGet(context) {
  try {
    // construct url to origin server
    const api = context.env.AZURE_COSMOS_URL + "/users";

    // TODO is the Cf-Access header included?
    const response = await fetch(api);

    return response;

  } catch (err) {
    return new Response('Error parsing JSON content', { status: 400 });
  }
}
