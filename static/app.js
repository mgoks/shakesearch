const Controller = {
  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}`).then((response) => {
      response.json().then((results) => {
        Controller.updateTable(results);
      });
    });
  },

  updateTable: (results) => {
    const table = document.getElementById("table-body");
    const rows = [];
    for (let result of results) {
      //const row = `<tr>${result}<tr/>`;
			const row = `
					<tr class="bg-white border-b dark:bg-gray-900 dark:border-gray-700">
						<td class="py-4 px-6">
							${result}
						</td>
					</tr>`
			rows.push(row);
    }
    table.innerHTML = rows.join('\n');
  },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
