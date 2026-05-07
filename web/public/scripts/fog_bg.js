(function () {
  if (!VANTA) return null
  setTimeout(() => {
    VANTA.FOG({
      el: "#body",
      mouseControls: true,
      touchControls: true,
      gyroControls: false,
      minHeight: 200.00,
      minWidth: 200.00,
      highlightColor: 0x9b7200,
      midtoneColor: 0x8e0a00,
      lowlightColor: 0xffffff,
      baseColor: 0x0
    })
  })
})()
