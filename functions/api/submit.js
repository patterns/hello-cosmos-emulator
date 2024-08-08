/**
 * POST /api/submit
 */
export async function onRequestPost(context) {
  try {
    let input = await context.request.formData();
    // Convert FormData to JSON
    // NOTE: Allows multiple values per key
    let output = {};
/**************************
    for (let [key, value] of input) {
      let tmp = output[key];
      if (tmp === undefined) {
        output[key] = value;
      } else {
        output[key] = [].concat(tmp, value);
      }
    }*/

    output["title"] = input["title"];
    output["description"] = input["description"];
    output["url"] = input["url"];

    let data = JSON.stringify(output);
    // submit input to our mongo backed api
    const api = context.env.AZURE_COSMOS_URL + "/courses";

    // TODO is the Cf-Access header included?
    const res = await fetch(api, {
        method: "POST",
        body: data,
    });
    return res;


  } catch (err) {
    return new Response('Error parsing JSON content', { status: 400 });
  }
}
