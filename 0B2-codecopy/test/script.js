document.addEventListener('DOMContentLoaded', function() {
  const copyButtons = document.querySelectorAll('.copy-button');
  console.log(copyButtons);

  copyButtons.forEach(button => {
    console.log('button', button);
    button.addEventListener('click', function() {
      console.log('clicked', this);
      const codeBlock = this.nextElementSibling;
      const codeText = codeBlock.querySelector('code').innerText;

      const tempInput = document.createElement('textarea');
      tempInput.value = codeText;
      document.body.appendChild(tempInput);

      tempInput.select();
      document.execCommand('copy');

      document.body.removeChild(tempInput);

      this.innerHTML='<i>Copied</i>';
      setTimeout(() => {
        this.innerHTML='Copy code';
      }, 2000); // Reset button state after 2 seconds
    });
  });
});
