Informe DEVOPS/SEGURIDAD — PROPOSE fase 2, revisión documental independiente.

**JUICIO:** diseño viable para continuar PROPOSE; no habilita APPLY. No identifico bloqueantes adicionales para esta entrega exclusivamente documental. Los siguientes contratos deben resolverse antes de implementar sus respectivos changes:

- **Firma/confianza (§8):** la raíz externa confiable, identidad autorizada y revocación siguen `(TBD)`. Es un gate correctamente declarado. Concretar cómo se fija configuración/entorno Git confiables y probar que configuración local maliciosa no selecciona ejecutables ni altera la política. Una firma válida no demuestra que una persona comprendió el alcance; el diseño evita equipararla automáticamente con consentimiento.
- **Hooks (§9):** falta explicitar el contrato del comando instalado: ruta del ejecutable, quoting por plataforma, tratamiento de eventos como datos y comportamiento cuando falla el capturador. Exigir que nombres de change/rol y contenido del evento nunca se interpolen como shell. Preservar orden y semántica de hooks personales, además de conservar sus bytes.
- **Paths y roles (§11):** extender explícitamente a instalación, actualización y retirada de roles las defensas contra traversal, symlinks y cambios concurrentes ya exigidas para captura. Validar ownership contra contenido actual; un manifiesto propio no autoriza sobrescribir un archivo posteriormente modificado. No declarar aislamiento fuerte mediante permisos cuando constructor y capturador comparten usuario.
- **Backup/recuperación (§§10–11):** son acertados snapshot previo, fallo cerrado y recuperación interna del paso fallido. Precisar que restaurar también compara el estado actual para evitar eliminar modificaciones personales concurrentes. La alternativa documentada de recuperación incompleta debe permanecer visible y bloquear éxito.

**EVIDENCIA:** la propuesta distingue VERIFIED/WRITTEN_UNVERIFIED, rechaza firmas desconocidas, evita programas de verificación arbitrarios del repo, conserva originales, exige instalación por adapters, protege archivos personales y reconoce que pipeline no revierte el paso fallido. El reemplazo destructivo de `bash` requiere aprobación específica J6 y compatibilidad demostrada.

**SUPOSICIONES:** los futuros fixtures reproducirán los contratos reales de cada versión y plataforma. No se acredita hoy eficacia de firma, hooks, backup ni recuperación.

**Fuentes y límites:** lectura de `AGENTS.md`, exploración y propuesta del 2026-10-23, especialmente §§8–11 y tabla de pruebas. Sin otros informes, escrituras, tests/build, configuración global ni `.env`. No verifiqué documentación externa ni implementación. Este informe no constituye aprobación humana.
