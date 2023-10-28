const imageFileInput = document.getElementById("images")! as HTMLInputElement;
const qualityInput = document.getElementById("quality") as HTMLInputElement;

// When you soft refresh the input is still populated
onQualityChange();
if (imageFileInput.files?.length! > 0) {
  onInputImageChange();
}

imageFileInput.addEventListener("change", onInputImageChange);
qualityInput.addEventListener("input", onQualityChange);

function onInputImageChange() {
  const filesUL = document.getElementById("files")! as HTMLUListElement;
  for (const file of imageFileInput.files!) {
    const li = document.createElement("li");
    li.textContent = `${file.name} - ${kbToMB(file.size)} ${file.type}`;
    filesUL.appendChild(li);
  }
  console.log(imageFileInput.files);
}

function onQualityChange() {
  const qualityIndicator = document.getElementById(
    "quality-indicator"
  )! as HTMLSpanElement;
  qualityIndicator.innerText = qualityInput.value;
}

function kbToMB(n: number): string {
  let ret: number;
  let suffix: string;
  if (n < 1e5) {
    ret = n / 1e3;
    suffix = "KB";
  } else {
    ret = n / 1e5;
    suffix = "MB";
  }
  return ret.toFixed(2) + suffix;
}
