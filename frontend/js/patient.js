import { apiListMedicalData, apiDeleteMedicalData } from "./api.js";

const user = JSON.parse(localStorage.getItem("user"));
if (!user) {
    location.href = "../index.html";
}

document.getElementById("patientName").textContent = `${user.firstName} ${user.lastName}`;
document.getElementById("patientMeta").textContent = user.email;

async function loadDiagnoses() {
    const container = document.getElementById("diagnosesList");
    try {
        const res = await apiListMedicalData(user.id);
        const list = res.medicalData || res || [];
        container.innerHTML = list.length
            ? list
                .map(
                    (d) => `
          <div class="diagnosis-item">
            <div>
              <h4>${d.diagnosis || "Без названия"}</h4>
              <p class="diagnosis-meta">Date: ${new Date(d.createdAt).toLocaleDateString()}</p>
            </div>
            <button class="link-btn" data-id="${d.id}">🗑</button>
          </div>`
                )
                .join("")
            : `<p class="text-gray">Нет диагнозов</p>`;

        container.querySelectorAll(".link-btn").forEach((btn) =>
            btn.addEventListener("click", async () => {
                await apiDeleteMedicalData(btn.dataset.id);
                loadDiagnoses();
            })
        );
    } catch (err) {
        container.innerHTML = `<p class="text-gray">Ошибка загрузки: ${err.message}</p>`;
    }
}

document.getElementById("refreshBtn").onclick = loadDiagnoses;
loadDiagnoses();
