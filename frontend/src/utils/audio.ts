export const playNotificationSound = (volume: number) => {
  try {
    const AudioContextClass = window.AudioContext || (window as any).webkitAudioContext;
    if (!AudioContextClass) return;
    const ctx = new AudioContextClass();
    
    const osc1 = ctx.createOscillator();
    const osc2 = ctx.createOscillator();
    const gainNode = ctx.createGain();
    
    osc1.type = "sine";
    osc2.type = "sine";
    
    osc1.frequency.setValueAtTime(587.33, ctx.currentTime);
    osc2.frequency.setValueAtTime(880.00, ctx.currentTime);
    
    gainNode.gain.setValueAtTime(0, ctx.currentTime);
    gainNode.gain.linearRampToValueAtTime(volume, ctx.currentTime + 0.05);
    gainNode.gain.exponentialRampToValueAtTime(0.0001, ctx.currentTime + 0.5);
    
    osc1.connect(gainNode);
    osc2.connect(gainNode);
    gainNode.connect(ctx.destination);
    
    osc1.start(ctx.currentTime);
    osc2.start(ctx.currentTime);
    
    osc1.stop(ctx.currentTime + 0.6);
    osc2.stop(ctx.currentTime + 0.6);

    setTimeout(() => {
      ctx.close().catch((err: any) => {
        console.error("Failed to close AudioContext", err);
      });
    }, 700);
  } catch (e) {
    console.error("Failed to play synthesized notification sound", e);
  }
};
