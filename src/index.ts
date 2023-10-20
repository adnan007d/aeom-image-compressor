const imageFileInput = document.getElementById("images")! as HTMLInputElement;

imageFileInput.addEventListener("change", onInputImageChange);

function onInputImageChange(_e: Event) {
  const filesUL = document.getElementById("files")! as HTMLUListElement;
  for (const file of imageFileInput.files!) {
    const li = document.createElement("li");
    li.textContent = `${file.name} - ${file.size} ${file.type}`;
    filesUL.appendChild(li);
  }
  console.log(imageFileInput.files);
}
