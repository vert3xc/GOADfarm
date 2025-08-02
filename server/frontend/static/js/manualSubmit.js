const apiToken = window.API_KEY;
        let currentPage = 1;
        const perPage = 20;

        function extractFlags(text) {
            const regex = window.FLAG_FORMAT;
            return [...text.matchAll(regex)].map(match => ({
                flag: match[0],
                sploit: "Manual",
                team: "*"
            }));
        }

        async function submitFlags() {
            const text = document.getElementById("flag-input").value;
            const flags = extractFlags(text);
            if (flags.length === 0) {
                alert("No flags found.");
                return;
            }

            await fetch("/api/post_flags", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-API-Token": apiToken
                },
                body: JSON.stringify(flags)
            });

            alert("Submitted " + flags.length + " flags!");
            loadFlags();
        }

        async function loadFlags(page = 1) {
            const params = new URLSearchParams();
            params.append("page", page);
            params.append("limit", perPage);

            const filters = {
                sploit: document.getElementById("filter-sploit").value,
                team: document.getElementById("filter-team").value,
                status: document.getElementById("filter-status").value,
                response: document.getElementById("filter-response").value,
            };

            for (const key in filters) {
                if (filters[key]) {
                    params.append(key, filters[key]);
                }
            }

            const res = await fetch(`/api/list_flags?${params.toString()}`, {
                headers: { "X-API-Token": apiToken }
            });
            const data = await res.json();

            const tbody = document.querySelector("#flag-table tbody");
            tbody.innerHTML = "";
            console.log(data.flags);
            console.log(data.total_pages);
            for (const flag of data.flags) {
                const row = document.createElement("tr");
                const tdTime = document.createElement("td");
                tdTime.textContent = new Date(flag.Time).toLocaleTimeString();
                row.appendChild(tdTime);
                const tdFlag = document.createElement("td");
                tdFlag.textContent = flag.Flag;
                row.appendChild(tdFlag);
                const tdSploit = document.createElement("td");
                tdSploit.textContent = flag.Sploit;
                row.appendChild(tdSploit);
                const tdTeam = document.createElement("td");
                tdTeam.textContent = flag.Team;
                row.appendChild(tdTeam);
                const tdStatus = document.createElement("td");
                tdStatus.textContent = flag.Status;
                row.appendChild(tdStatus);
                const tdResponse = document.createElement("td");
                tdResponse.textContent = flag.Response;
                row.appendChild(tdResponse);
                tbody.appendChild(row);
            }
            const pagination = document.getElementById("pagination");
            pagination.innerHTML = "";
            for (let i = 1; i <= data.total_pages; i++) {
                const btn = document.createElement("button");
                btn.textContent = i;
                if (i === page) btn.disabled = true;
                btn.onclick = () => { currentPage = i; loadFlags(i); };
                pagination.appendChild(btn);
            }
        }

        loadFlags();