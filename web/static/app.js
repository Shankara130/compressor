function upload() {
  const file = document.getElementById("file").files[0]
  const form = new FormData()
  form.append("file", file)

  fetch("/upload", { method: "POST", body: form })
    .then(r => r.json())
    .then(d => poll(d.job_id))
}

function poll(id) {
  const i = setInterval(() => {
    fetch("/status/" + id)
      .then(r => r.json())
      .then(j => {
        document.getElementById("status").innerText =
          j.status + " " + j.progress + "%"

        if (j.status === "DONE") {
          clearInterval(i)
          window.location = "/download/" + id
        }
      })
  }, 1000)
}
