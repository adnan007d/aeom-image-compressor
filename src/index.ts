const imageFileInput = document.getElementById("images")! as HTMLInputElement;

// When you soft refresh the input is still populated
if (imageFileInput.files?.length! > 0) {
  onInputImageChange();
}

imageFileInput.addEventListener("change", onInputImageChange);

function onInputImageChange() {
  const filesUL = document.getElementById("files")! as HTMLUListElement;
  for (const file of imageFileInput.files!) {
    const li = document.createElement("li");
    li.textContent = `${file.name} - ${file.size} ${file.type}`;
    filesUL.appendChild(li);
  }
  console.log(imageFileInput.files);
}
