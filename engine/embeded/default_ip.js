function get_ip_by_cf() {
	const urls = ["https://1.1.1.1/cdn-cgi/trace", "https://[2606:4700:4700::1111]/cdn-cgi/trace"];
	const ipret = [];
	urls.forEach((url) => {
		const cf_content = (get(fetch(url, {
			headers: {
				'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36',
			},
			retry: 1,
			timeout: 3000,
		}), "body") || "").trim();

		const ip = (cf_content.match(/ip=(\S+)/)?.[1] || '').trim();
		if (ip) ipret.push(ip);
	});
	return ipret;
}
const ip_resolve_default = get_ip_by_cf;
