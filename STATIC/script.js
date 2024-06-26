const navLinks = document.querySelectorAll('nav ul li a');

navLinks.forEach(link => {
  link.addEventListener('mouseover', () => {
    const submenu = link.nextElementSibling;
    submenu.style.display = 'block';
  });

  link.addEventListener('mouseout', () => {
    const submenu = link.nextElementSibling;
    submenu.style.display = 'none';
  });
});