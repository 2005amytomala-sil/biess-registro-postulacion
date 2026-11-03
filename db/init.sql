-- ============================================================================
-- BLOQUE 1: DDL Normalizado (3FN)
-- Sistema de Registro y Postulación al Concurso Público de Méritos del BIESS
-- ============================================================================

-- Habilitar extensiones necesarias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- CATÁLOGOS: Paises, Provincias, Cantones
-- ============================================================================

CREATE TABLE paises (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE provincias (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cantones (
    id SERIAL PRIMARY KEY,
    provincia_id INTEGER NOT NULL REFERENCES provincias(id) ON DELETE RESTRICT,
    nombre VARCHAR(100) NOT NULL,
    UNIQUE (provincia_id, nombre),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para catálogos geográficos
CREATE INDEX idx_cantones_provincia_id ON cantones(provincia_id);

-- ============================================================================
-- TABLA: postulantes
-- ============================================================================
CREATE TABLE postulantes (
    id SERIAL PRIMARY KEY,
    cedula VARCHAR(20) NOT NULL UNIQUE,
    correo VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para postulantes
CREATE INDEX idx_postulantes_cedula ON postulantes(cedula);
CREATE INDEX idx_postulantes_correo ON postulantes(correo);

-- ============================================================================
-- TABLA: postulaciones
-- ============================================================================
CREATE TABLE postulaciones (
    id SERIAL PRIMARY KEY,
    postulante_id INTEGER NOT NULL UNIQUE REFERENCES postulantes(id) ON DELETE CASCADE,
    codigo VARCHAR(50) NOT NULL UNIQUE,
    sector VARCHAR(20) NOT NULL CHECK (sector IN ('afiliados', 'jubilados', 'empleadores')),
    estado VARCHAR(20) NOT NULL DEFAULT 'borrador' CHECK (estado IN ('borrador', 'enviada', 'cerrada')),
    motivo_cierre TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para postulaciones
CREATE INDEX idx_postulaciones_postulante_id ON postulaciones(postulante_id);
CREATE INDEX idx_postulaciones_codigo ON postulaciones(codigo);
CREATE INDEX idx_postulaciones_estado ON postulaciones(estado);

-- Trigger para actualizar updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_postulaciones_updated_at
    BEFORE UPDATE ON postulaciones
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- TABLA: datos_personales
-- ============================================================================
CREATE TABLE datos_personales (
    postulacion_id INTEGER PRIMARY KEY REFERENCES postulaciones(id) ON DELETE CASCADE,
    nombres VARCHAR(200) NOT NULL,
    apellidos VARCHAR(200) NOT NULL,
    cedula VARCHAR(20) NOT NULL,
    genero VARCHAR(20) NOT NULL CHECK (genero IN ('masculino', 'femenino', 'otro')),
    fecha_nacimiento DATE NOT NULL,
    nacionalidad VARCHAR(50) NOT NULL,
    estado_civil VARCHAR(50) NOT NULL,
    pareja_nombres VARCHAR(200),
    pareja_cedula VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- TABLA: domicilio_contacto
-- ============================================================================
CREATE TABLE domicilio_contacto (
    postulacion_id INTEGER PRIMARY KEY REFERENCES postulaciones(id) ON DELETE CASCADE,
    pais_id INTEGER NOT NULL REFERENCES paises(id) ON DELETE RESTRICT,
    provincia_id INTEGER REFERENCES provincias(id) ON DELETE RESTRICT,
    canton_id INTEGER REFERENCES cantones(id) ON DELETE RESTRICT,
    provincia_texto VARCHAR(100),
    ciudad_texto VARCHAR(100),
    celular VARCHAR(20) NOT NULL,
    correo_1 VARCHAR(255) NOT NULL,
    correo_2 VARCHAR(255),
    calle_principal VARCHAR(200),
    numero VARCHAR(50),
    calle_secundaria VARCHAR(200),
    sector VARCHAR(200),
    referencia TEXT,
    telefono_domicilio VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (pais_id = 1 AND provincia_id IS NOT NULL AND canton_id IS NOT NULL) OR
        (pais_id != 1 AND provincia_id IS NULL AND canton_id IS NULL)
    )
);

-- Índices para domicilio_contacto
CREATE INDEX idx_domicilio_pais_id ON domicilio_contacto(pais_id);
CREATE INDEX idx_domicilio_provincia_id ON domicilio_contacto(provincia_id);
CREATE INDEX idx_domicilio_canton_id ON domicilio_contacto(canton_id);

-- ============================================================================
-- TABLA: familiares_iess
-- ============================================================================
CREATE TABLE familiares_iess (
    id SERIAL PRIMARY KEY,
    postulacion_id INTEGER NOT NULL REFERENCES postulaciones(id) ON DELETE CASCADE,
    nombres VARCHAR(200) NOT NULL,
    parentesco VARCHAR(100) NOT NULL,
    institucion VARCHAR(20) NOT NULL CHECK (institucion IN ('IESS', 'BIESS')),
    area VARCHAR(200),
    cargo VARCHAR(200),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para familiares_iess
CREATE INDEX idx_familiares_postulacion_id ON familiares_iess(postulacion_id);

-- ============================================================================
-- TABLA: respaldo_organizacion
-- ============================================================================
CREATE TABLE respaldo_organizacion (
    postulacion_id INTEGER PRIMARY KEY REFERENCES postulaciones(id) ON DELETE CASCADE,
    tiene_respaldo BOOLEAN NOT NULL DEFAULT FALSE,
    nombre_organizacion VARCHAR(300),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        (tiene_respaldo = TRUE AND nombre_organizacion IS NOT NULL) OR
        (tiene_respaldo = FALSE)
    )
);

-- ============================================================================
-- TABLA: formacion_academica
-- ============================================================================
CREATE TABLE formacion_academica (
    id SERIAL PRIMARY KEY,
    postulacion_id INTEGER NOT NULL REFERENCES postulaciones(id) ON DELETE CASCADE,
    nivel VARCHAR(100) NOT NULL,
    institucion VARCHAR(300) NOT NULL,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('Nacional', 'Extranjera')),
    titulo VARCHAR(300) NOT NULL,
    area_estudios VARCHAR(200),
    registro_senescyt VARCHAR(100),
    fecha_registro_senescyt DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices para formacion_academica
CREATE INDEX idx_formacion_postulacion_id ON formacion_academica(postulacion_id);

-- ============================================================================
-- TABLA: experiencia_profesional
-- ============================================================================
CREATE TABLE experiencia_profesional (
    id SERIAL PRIMARY KEY,
    postulacion_id INTEGER NOT NULL REFERENCES postulaciones(id) ON DELETE CASCADE,
    tipo VARCHAR(20) NOT NULL CHECK (tipo IN ('General', 'Direccion')),
    cargo VARCHAR(200) NOT NULL,
    institucion VARCHAR(300) NOT NULL,
    fecha_inicio DATE NOT NULL,
    fecha_fin DATE,
    actividades_relacionadas TEXT,
    descripcion VARCHAR(1000),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CHECK (fecha_fin IS NULL OR fecha_fin >= fecha_inicio)
);

-- Índices para experiencia_profesional
CREATE INDEX idx_experiencia_postulacion_id ON experiencia_profesional(postulacion_id);

-- ============================================================================
-- TABLA: inhabilidades
-- ============================================================================
CREATE TABLE inhabilidades (
    postulacion_id INTEGER NOT NULL REFERENCES postulaciones(id) ON DELETE CASCADE,
    codigo_pregunta VARCHAR(20) NOT NULL CHECK (codigo_pregunta IN (
        '1.9.1', '1.9.2', '1.9.3', '1.9.4', '1.9.5',
        '1.9.6', '1.9.7', '1.9.8', '1.9.9', '1.9.10',
        '1.9.11', '1.9.12', '1.9.13'
    )),
    respuesta BOOLEAN NOT NULL,
    detalle JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (postulacion_id, codigo_pregunta),
    CHECK (
        (codigo_pregunta = '1.9.7' AND detalle IS NOT NULL) OR
        (codigo_pregunta != '1.9.7')
    )
);

-- Índices para inhabilidades
CREATE INDEX idx_inhabilidades_postulacion_id ON inhabilidades(postulacion_id);
CREATE INDEX idx_inhabilidades_codigo_pregunta ON inhabilidades(codigo_pregunta);

-- ============================================================================
-- FIN DEL BLOQUE 1: DDL
-- ============================================================================

-- ============================================================================
-- BLOQUE 2: Semilla de Países y Provincias del Ecuador
-- ============================================================================

-- ----------------------------------------------------------------------------
-- INSERT: Paises (Ecuador con id=1 para el CHECK de cascada geográfica)
-- ----------------------------------------------------------------------------
INSERT INTO paises (id, nombre) VALUES
(1, 'Ecuador'),
(2, 'Afganistán'),
(3, 'Albania'),
(4, 'Alemania'),
(5, 'Andorra'),
(6, 'Angola'),
(7, 'Antigua y Barbuda'),
(8, 'Arabia Saudita'),
(9, 'Argelia'),
(10, 'Argentina'),
(11, 'Armenia'),
(12, 'Australia'),
(13, 'Austria'),
(14, 'Azerbaiyán'),
(15, 'Bahamas'),
(16, 'Bangladés'),
(17, 'Barbados'),
(18, 'Baréin'),
(19, 'Bélgica'),
(20, 'Belice'),
(21, 'Benín'),
(22, 'Bielorrusia'),
(23, 'Birmania'),
(24, 'Bolivia'),
(25, 'Bosnia y Herzegovina'),
(26, 'Botsuana'),
(27, 'Brasil'),
(28, 'Brunéi'),
(29, 'Bulgaria'),
(30, 'Burkina Faso'),
(31, 'Burundi'),
(32, 'Bután'),
(33, 'Cabo Verde'),
(34, 'Camboya'),
(35, 'Camerún'),
(36, 'Canadá'),
(37, 'Catar'),
(38, 'Chad'),
(39, 'Chile'),
(40, 'China'),
(41, 'Chipre'),
(42, 'Ciudad del Vaticano'),
(43, 'Colombia'),
(44, 'Comoras'),
(45, 'Corea del Norte'),
(46, 'Corea del Sur'),
(47, 'Costa de Marfil'),
(48, 'Costa Rica'),
(49, 'Croacia'),
(50, 'Cuba'),
(51, 'Dinamarca'),
(52, 'Dominica'),
(53, 'Ecuador'),
(54, 'Egipto'),
(55, 'El Salvador'),
(56, 'Emiratos Árabes Unidos'),
(57, 'Eritrea'),
(58, 'Eslovaquia'),
(59, 'Eslovenia'),
(60, 'España'),
(61, 'Estados Unidos'),
(62, 'Estonia'),
(63, 'Etiopía'),
(64, 'Filipinas'),
(65, 'Finlandia'),
(66, 'Fiyi'),
(67, 'Francia'),
(68, 'Gabón'),
(69, 'Gambia'),
(70, 'Georgia'),
(71, 'Ghana'),
(72, 'Granada'),
(73, 'Grecia'),
(74, 'Guatemala'),
(75, 'Guinea'),
(76, 'Guinea Ecuatorial'),
(77, 'Guinea-Bisáu'),
(78, 'Guyana'),
(79, 'Haití'),
(80, 'Honduras'),
(81, 'Hungría'),
(82, 'India'),
(83, 'Indonesia'),
(84, 'Irak'),
(85, 'Irán'),
(86, 'Irlanda'),
(87, 'Islandia'),
(88, 'Islas Marshall'),
(89, 'Islas Salomón'),
(90, 'Israel'),
(91, 'Italia'),
(92, 'Jamaica'),
(93, 'Japón'),
(94, 'Jordania'),
(95, 'Kazajistán'),
(96, 'Kenia'),
(97, 'Kirguistán'),
(98, 'Kiribati'),
(99, 'Kuwait'),
(100, 'Laos'),
(101, 'Lesoto'),
(102, 'Letonia'),
(103, 'Líbano'),
(104, 'Liberia'),
(105, 'Libia'),
(106, 'Liechtenstein'),
(107, 'Lituania'),
(108, 'Luxemburgo'),
(109, 'Macedonia del Norte'),
(110, 'Madagascar'),
(111, 'Malasia'),
(112, 'Malaui'),
(113, 'Maldivas'),
(114, 'Malí'),
(115, 'Malta'),
(116, 'Marruecos'),
(117, 'Mauricio'),
(118, 'Mauritania'),
(119, 'México'),
(120, 'Micronesia'),
(121, 'Moldavia'),
(122, 'Mónaco'),
(123, 'Mongolia'),
(124, 'Montenegro'),
(125, 'Mozambique'),
(126, 'Namibia'),
(127, 'Nauru'),
(128, 'Nepal'),
(129, 'Nicaragua'),
(130, 'Níger'),
(131, 'Nigeria'),
(132, 'Noruega'),
(133, 'Nueva Zelanda'),
(134, 'Omán'),
(135, 'Países Bajos'),
(136, 'Pakistán'),
(137, 'Palaos'),
(138, 'Panamá'),
(139, 'Papúa Nueva Guinea'),
(140, 'Paraguay'),
(141, 'Perú'),
(142, 'Polonia'),
(143, 'Portugal'),
(144, 'Reino Unido'),
(145, 'República Centroafricana'),
(146, 'República Checa'),
(147, 'República del Congo'),
(148, 'República Democrática del Congo'),
(149, 'República Dominicana'),
(150, 'Ruanda'),
(151, 'Rumania'),
(152, 'Rusia'),
(153, 'Samoa'),
(154, 'San Cristóbal y Nieves'),
(155, 'San Marino'),
(156, 'San Vicente y las Granadinas'),
(157, 'Santa Lucía'),
(158, 'Santo Tomé y Príncipe'),
(159, 'Senegal'),
(160, 'Serbia'),
(161, 'Seychelles'),
(162, 'Sierra Leona'),
(163, 'Singapur'),
(164, 'Siria'),
(165, 'Somalia'),
(166, 'Sri Lanka'),
(167, 'Suazilandia'),
(168, 'Sudáfrica'),
(169, 'Sudán'),
(170, 'Sudán del Sur'),
(171, 'Suecia'),
(172, 'Suiza'),
(173, 'Surinam'),
(174, 'Tailandia'),
(175, 'Tanzania'),
(176, 'Tayikistán'),
(177, 'Timor Oriental'),
(178, 'Togo'),
(179, 'Tonga'),
(180, 'Trinidad y Tobago'),
(181, 'Túnez'),
(182, 'Turkmenistán'),
(183, 'Turquía'),
(184, 'Tuvalu'),
(185, 'Ucrania'),
(186, 'Uganda'),
(187, 'Uruguay'),
(188, 'Uzbekistán'),
(189, 'Vanuatu'),
(190, 'Venezuela'),
(191, 'Vietnam'),
(192, 'Yemen'),
(193, 'Yibuti'),
(194, 'Zambia'),
(195, 'Zimbabue')
ON CONFLICT (nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- INSERT: Provincias del Ecuador (24 provincias)
-- ----------------------------------------------------------------------------
INSERT INTO provincias (id, nombre) VALUES
(1, 'Azuay'),
(2, 'Bolívar'),
(3, 'Cañar'),
(4, 'Carchi'),
(5, 'Chimborazo'),
(6, 'Cotopaxi'),
(7, 'El Oro'),
(8, 'Esmeraldas'),
(9, 'Galápagos'),
(10, 'Guayas'),
(11, 'Imbabura'),
(12, 'Loja'),
(13, 'Los Ríos'),
(14, 'Manabí'),
(15, 'Morona Santiago'),
(16, 'Napo'),
(17, 'Orellana'),
(18, 'Pastaza'),
(19, 'Pichincha'),
(20, 'Santa Elena'),
(21, 'Santo Domingo de los Tsáchilas'),
(22, 'Sucumbíos'),
(23, 'Tungurahua'),
(24, 'Zamora Chinchipe')
ON CONFLICT (nombre) DO NOTHING;

-- ============================================================================
-- FIN DEL BLOQUE 2: Países y Provincias
-- ============================================================================

-- ============================================================================
-- BLOQUE 3: Cantones por Provincia (según requerimiento.pdf)
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Azuay (provincia_id=1)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(1, 'Cuenca'),
(1, 'Gualaceo'),
(1, 'Chordeleg'),
(1, 'Sígsig'),
(1, 'Girón'),
(1, 'Santa Isabel'),
(1, 'Pucará'),
(1, 'Camilo Ponce Enríquez')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Bolívar (provincia_id=2)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(2, 'Guaranda'),
(2, 'Chillanes'),
(2, 'Chimbo'),
(2, 'Echeandía'),
(2, 'San Miguel'),
(2, 'Caluma'),
(2, 'Las Naves')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Cañar (provincia_id=3)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(3, 'Azogues'),
(3, 'Biblián'),
(3, 'Cañar'),
(3, 'Déleg'),
(3, 'El Tambo'),
(3, 'Suscal'),
(3, 'La Troncal')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Carchi (provincia_id=4)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(4, 'Tulcán'),
(4, 'Bolívar'),
(4, 'Espejo'),
(4, 'Mira'),
(4, 'Montúfar'),
(4, 'San Pedro de Huaca')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Chimborazo (provincia_id=5)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(5, 'Riobamba'),
(5, 'Alausí'),
(5, 'Colta'),
(5, 'Chunchi'),
(5, 'Guamote'),
(5, 'Guano'),
(5, 'Pallatanga'),
(5, 'Penipe'),
(5, 'Cumandá')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Cotopaxi (provincia_id=6)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(6, 'Latacunga'),
(6, 'La Maná'),
(6, 'Pangua'),
(6, 'Pujilí'),
(6, 'Salcedo'),
(6, 'Saquisilí'),
(6, 'Sigchos')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- El Oro (provincia_id=7)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(7, 'Machala'),
(7, 'Arenillas'),
(7, 'Atahualpa'),
(7, 'Balsas'),
(7, 'Chilla'),
(7, 'El Guabo'),
(7, 'Huaquillas'),
(7, 'Las Lajas'),
(7, 'Marcabelí'),
(7, 'Pasaje'),
(7, 'Piñas'),
(7, 'Portovelo'),
(7, 'Santa Rosa'),
(7, 'Zaruma')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Esmeraldas (provincia_id=8)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(8, 'Esmeraldas'),
(8, 'Atacames'),
(8, 'Eloy Alfaro'),
(8, 'Muisne'),
(8, 'Quinindé'),
(8, 'Rioverde'),
(8, 'San Lorenzo')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Galápagos (provincia_id=9)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(9, 'San Cristóbal'),
(9, 'Isabela'),
(9, 'Santa Cruz')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Guayas (provincia_id=10)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(10, 'Guayaquil'),
(10, 'Durán'),
(10, 'Milagro'),
(10, 'Samborondón'),
(10, 'Daule'),
(10, 'Playas'),
(10, 'Balzar'),
(10, 'Colimes'),
(10, 'El Empalme'),
(10, 'El Triunfo'),
(10, 'Isidro Ayora'),
(10, 'Lomas de Sargentillo'),
(10, 'Naranjal'),
(10, 'Naranjito'),
(10, 'Nobol'),
(10, 'Pedro Carbo'),
(10, 'Salitre'),
(10, 'Simón Bolívar'),
(10, 'Yaguachi')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Imbabura (provincia_id=11)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(11, 'Ibarra'),
(11, 'Antonio Ante'),
(11, 'Cotacachi'),
(11, 'Otavalo'),
(11, 'Pimampiro'),
(11, 'San Miguel de Urcuquí')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Loja (provincia_id=12)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(12, 'Loja'),
(12, 'Calvas'),
(12, 'Catamayo'),
(12, 'Celica'),
(12, 'Chaguarpamba'),
(12, 'Espíndola'),
(12, 'Gonzanamá'),
(12, 'Macará'),
(12, 'Paltas'),
(12, 'Puyango'),
(12, 'Saraguro'),
(12, 'Sozoranga'),
(12, 'Zapotillo')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Los Ríos (provincia_id=13)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(13, 'Babahoyo'),
(13, 'Baba'),
(13, 'Buena Fe'),
(13, 'Mocache'),
(13, 'Montalvo'),
(13, 'Palenque'),
(13, 'Puebloviejo'),
(13, 'Quevedo'),
(13, 'Urdaneta'),
(13, 'Ventanas'),
(13, 'Vinces'),
(13, 'Valencia')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Manabí (provincia_id=14)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(14, 'Portoviejo'),
(14, 'Manta'),
(14, 'Chone'),
(14, 'Bolívar'),
(14, 'El Carmen'),
(14, 'Flavio Alfaro'),
(14, 'Jipijapa'),
(14, 'Junín'),
(14, 'Montecristi'),
(14, 'Paján'),
(14, 'Pichincha'),
(14, 'Rocafuerte'),
(14, 'Santa Ana'),
(14, 'Sucre'),
(14, 'Tosagua'),
(14, '24 de Mayo'),
(14, 'Pedernales'),
(14, 'Olmedo'),
(14, 'Puerto López'),
(14, 'Jama'),
(14, 'Jaramijó'),
(14, 'San Vicente')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Morona Santiago (provincia_id=15)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(15, 'Morona'),
(15, 'Gualaquiza'),
(15, 'Limón Indanza'),
(15, 'Palora'),
(15, 'San Juan Bosco'),
(15, 'Santiago'),
(15, 'Sucúa'),
(15, 'Taisha'),
(15, 'Tiwintza')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Napo (provincia_id=16)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(16, 'Tena'),
(16, 'Archidona'),
(16, 'Carlos Julio Arosemena Tola'),
(16, 'El Chaco'),
(16, 'Quijos')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Orellana (provincia_id=17)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(17, 'Francisco de Orellana'),
(17, 'Aguarico'),
(17, 'La Joya de los Sachas'),
(17, 'Loreto')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Pastaza (provincia_id=18)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(18, 'Pastaza'),
(18, 'Arajuno'),
(18, 'Mera'),
(18, 'Santa Clara')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Pichincha (provincia_id=19)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(19, 'Quito'),
(19, 'Cayambe'),
(19, 'Mejía'),
(19, 'Pedro Moncayo'),
(19, 'Pedro Vicente Maldonado'),
(19, 'Puerto Quito'),
(19, 'Rumiñahui'),
(19, 'San Miguel de los Bancos')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Santa Elena (provincia_id=20)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(20, 'Santa Elena'),
(20, 'La Libertad')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Santo Domingo de los Tsáchilas (provincia_id=21)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(21, 'La Concordia'),
(21, 'Santo Domingo')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Sucumbíos (provincia_id=22)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(22, 'Cascales'),
(22, 'Cuyabeno'),
(22, 'Gonzalo Pizarro'),
(22, 'Lago Agrio'),
(22, 'Putumayo'),
(22, 'Shushufindi'),
(22, 'Sucumbíos')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Tungurahua (provincia_id=23)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(23, 'Ambato'),
(23, 'Baños de Agua Santa'),
(23, 'Cevallos'),
(23, 'Mocha'),
(23, 'Patate'),
(23, 'Pelileo'),
(23, 'Píllaro'),
(23, 'Quero'),
(23, 'Tisaleo')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ----------------------------------------------------------------------------
-- Zamora Chinchipe (provincia_id=24)
-- ----------------------------------------------------------------------------
INSERT INTO cantones (provincia_id, nombre) VALUES
(24, 'Centinela del Cóndor'),
(24, 'Chinchipe'),
(24, 'El Pangui'),
(24, 'Nangaritza'),
(24, 'Palanda'),
(24, 'Paquisha'),
(24, 'Yacuambi'),
(24, 'Yantzaza'),
(24, 'Zamora')
ON CONFLICT (provincia_id, nombre) DO NOTHING;

-- ============================================================================
-- FIN DEL BLOQUE 3: Cantones
-- ============================================================================
