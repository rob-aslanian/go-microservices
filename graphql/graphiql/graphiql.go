package graphiql

import "net/http"

func GetGraphiql() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	})
}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<!--link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.12.0/graphiql.css"/>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/2.0.3/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/16.2.0/umd/react.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/16.2.0/umd/react-dom.production.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.12.0/graphiql.min.js"></script-->

		<link rel="stylesheet" href="/libs/graphiql.css"/>
		<script src="/libs/fetch.min.js"></script>
		<script src="/libs/react.production.min.js"></script>
		<script src="/libs/react-dom.production.min.js"></script>
		<script src="/libs/graphiql.min.js"></script>
		<script src="/libs/fetcher.js"></script>
		<script src="/libs/ws.js"></script>
		<style>
		.cm-ws {
			visibility:hidden;
		}
		.CodeMirror pre{
			color: rgba(0,0,0,0);
		}
		.CodeMirror-cursors {
			color: rgba(0,0,0,0);
		}
		</style>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function fetchGQL(params) {
				return fetch("/graphql", {
					method: "post",
					body: JSON.stringify(params),
					credentials: "include",
				}).then(function (resp) {
					return resp.text();
				}).then(function (body) {
					try {
						return JSON.parse(body);
					} catch (error) {
						return body;
					}
				});
			}

			var subscriptionsClient = new window.SubscriptionsTransportWs.SubscriptionClient('ws://testing.lan/graphqlws', { reconnect: true });
			var subscriptionsFetcher = window.GraphiQLSubscriptionsFetcher.graphQLFetcher(subscriptionsClient, fetchGQL);
			ReactDOM.render(
							React.createElement(GraphiQL, {fetcher: subscriptionsFetcher}),
							document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
