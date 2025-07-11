/** @format */

import * as React from "react";
const CheckFailure = (props) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    xmlnsXlink="http://www.w3.org/1999/xlink"
    width={80}
    height={80}
    viewBox="0 0 256 256"
    xmlSpace="preserve"
    {...props}
  >
    <defs />
    <g
      style={{
        stroke: "none",
        strokeWidth: 0,
        strokeDasharray: "none",
        strokeLinecap: "butt",
        strokeLinejoin: "miter",
        strokeMiterlimit: 10,
        fill: "none",
        fillRule: "nonzero",
        opacity: 1,
      }}
      transform="translate(1.4065934065934016 1.4065934065934016) scale(2.81 2.81)"
    >
      <path
        d="M 45 90 C 20.187 90 0 69.813 0 45 C 0 20.187 20.187 0 45 0 c 24.813 0 45 20.187 45 45 C 90 69.813 69.813 90 45 90 z"
        style={{
          stroke: "none",
          strokeWidth: 1,
          strokeDasharray: "none",
          strokeLinecap: "butt",
          strokeLinejoin: "miter",
          strokeMiterlimit: 10,
          fill: "rgb(236,0,0)",
          fillRule: "nonzero",
          opacity: 1,
        }}
        transform=" matrix(1 0 0 1 0 0) "
        strokeLinecap="round"
      />
      <path
        d="M 28.5 65.5 c -1.024 0 -2.047 -0.391 -2.829 -1.172 c -1.562 -1.562 -1.562 -4.095 0 -5.656 l 33 -33 c 1.561 -1.562 4.096 -1.562 5.656 0 c 1.563 1.563 1.563 4.095 0 5.657 l -33 33 C 30.547 65.109 29.524 65.5 28.5 65.5 z"
        style={{
          stroke: "none",
          strokeWidth: 1,
          strokeDasharray: "none",
          strokeLinecap: "butt",
          strokeLinejoin: "miter",
          strokeMiterlimit: 10,
          fill: "rgb(255,255,255)",
          fillRule: "nonzero",
          opacity: 1,
        }}
        transform=" matrix(1 0 0 1 0 0) "
        strokeLinecap="round"
      />
      <path
        d="M 61.5 65.5 c -1.023 0 -2.048 -0.391 -2.828 -1.172 l -33 -33 c -1.562 -1.563 -1.562 -4.095 0 -5.657 c 1.563 -1.562 4.095 -1.562 5.657 0 l 33 33 c 1.563 1.562 1.563 4.095 0 5.656 C 63.548 65.109 62.523 65.5 61.5 65.5 z"
        style={{
          stroke: "none",
          strokeWidth: 1,
          strokeDasharray: "none",
          strokeLinecap: "butt",
          strokeLinejoin: "miter",
          strokeMiterlimit: 10,
          fill: "rgb(255,255,255)",
          fillRule: "nonzero",
          opacity: 1,
        }}
        transform=" matrix(1 0 0 1 0 0) "
        strokeLinecap="round"
      />
    </g>
  </svg>
);

const Menu = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" {...props}>
    <path
      fill="#fff"
      d="M0 96c0-17.7 14.3-32 32-32h384c17.7 0 32 14.3 32 32s-14.3 32-32 32H32c-17.7 0-32-14.3-32-32zm0 160c0-17.7 14.3-32 32-32h384c17.7 0 32 14.3 32 32s-14.3 32-32 32H32c-17.7 0-32-14.3-32-32zm448 160c0 17.7-14.3 32-32 32H32c-17.7 0-32-14.3-32-32s14.3-32 32-32h384c17.7 0 32 14.3 32 32z"
    />
  </svg>
);

const Buses = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" {...props}>
    <path
      fill="#ffffff"
      d="M224 0C348.8 0 448 35.2 448 80V96 416c0 17.7-14.3 32-32 32v32c0 17.7-14.3 32-32 32H352c-17.7 0-32-14.3-32-32V448H128v32c0 17.7-14.3 32-32 32H64c-17.7 0-32-14.3-32-32l0-32c-17.7 0-32-14.3-32-32V96 80C0 35.2 99.2 0 224 0zM64 128V256c0 17.7 14.3 32 32 32H352c17.7 0 32-14.3 32-32V128c0-17.7-14.3-32-32-32H96c-17.7 0-32 14.3-32 32zM80 400a32 32 0 1 0 0-64 32 32 0 1 0 0 64zm288 0a32 32 0 1 0 0-64 32 32 0 1 0 0 64z"
    />
  </svg>
);

export const Pickup = (props) => (
  <svg
    {...props}
    xmlns="http://www.w3.org/2000/svg"
    fill="none"
    viewBox="0 0 24 24"
    strokeWidth={1.5}
    stroke="currentColor"
  >
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M15 10.5a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
    />
    <path
      strokeLinecap="round"
      strokeLinejoin="round"
      d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1 1 15 0Z"
    />
  </svg>
);
const Route = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" {...props}>
    <path
      fill="#ffffff"
      d="M512 96c0 50.2-59.1 125.1-84.6 155c-3.8 4.4-9.4 6.1-14.5 5H320c-17.7 0-32 14.3-32 32s14.3 32 32 32h96c53 0 96 43 96 96s-43 96-96 96H139.6c8.7-9.9 19.3-22.6 30-36.8c6.3-8.4 12.8-17.6 19-27.2H416c17.7 0 32-14.3 32-32s-14.3-32-32-32H320c-53 0-96-43-96-96s43-96 96-96h39.8c-21-31.5-39.8-67.7-39.8-96c0-53 43-96 96-96s96 43 96 96zM117.1 489.1c-3.8 4.3-7.2 8.1-10.1 11.3l-1.8 2-.2-.2c-6 4.6-14.6 4-20-1.8C59.8 473 0 402.5 0 352c0-53 43-96 96-96s96 43 96 96c0 30-21.1 67-43.5 97.9c-10.7 14.7-21.7 28-30.8 38.5l-.6 .7zM128 352a32 32 0 1 0 -64 0 32 32 0 1 0 64 0zM416 128a32 32 0 1 0 0-64 32 32 0 1 0 0 64z"
    />
  </svg>
);

const Trip = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" {...props}>
    <path
      fill="#ffffff"
      d="M152 24c0-13.3-10.7-24-24-24s-24 10.7-24 24V64H64C28.7 64 0 92.7 0 128v16 48V448c0 35.3 28.7 64 64 64H384c35.3 0 64-28.7 64-64V192 144 128c0-35.3-28.7-64-64-64H344V24c0-13.3-10.7-24-24-24s-24 10.7-24 24V64H152V24zM48 192h80v56H48V192zm0 104h80v64H48V296zm128 0h96v64H176V296zm144 0h80v64H320V296zm80-48H320V192h80v56zm0 160v40c0 8.8-7.2 16-16 16H320V408h80zm-128 0v56H176V408h96zm-144 0v56H64c-8.8 0-16-7.2-16-16V408h80zM272 248H176V192h96v56z"
    />
  </svg>
);

const BusStation = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 384 512" {...props}>
    <path
      fill="#ffffff"
      d="M215.7 499.2C267 435 384 279.4 384 192C384 86 298 0 192 0S0 86 0 192c0 87.4 117 243 168.3 307.2c12.3 15.3 35.1 15.3 47.4 0zM192 128a64 64 0 1 1 0 128 64 64 0 1 1 0-128z"
    />
  </svg>
);

const Ticket = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 576 512" {...props}>
    <path
      fill="#ffffff"
      d="M64 64C28.7 64 0 92.7 0 128v64c0 8.8 7.4 15.7 15.7 18.6C34.5 217.1 48 235 48 256s-13.5 38.9-32.3 45.4C7.4 304.3 0 311.2 0 320v64c0 35.3 28.7 64 64 64H512c35.3 0 64-28.7 64-64V320c0-8.8-7.4-15.7-15.7-18.6C541.5 294.9 528 277 528 256s13.5-38.9 32.3-45.4c8.3-2.9 15.7-9.8 15.7-18.6V128c0-35.3-28.7-64-64-64H64zm64 112l0 160c0 8.8 7.2 16 16 16H432c8.8 0 16-7.2 16-16V176c0-8.8-7.2-16-16-16H144c-8.8 0-16 7.2-16 16zM96 160c0-17.7 14.3-32 32-32H448c17.7 0 32 14.3 32 32V352c0 17.7-14.3 32-32 32H128c-17.7 0-32-14.3-32-32V160z"
    />
  </svg>
);

const Employee = (props) => (
  <svg
    height="200px"
    width="200px"
    id="_x32_"
    xmlns="http://www.w3.org/2000/svg"
    xmlnsXlink="http://www.w3.org/1999/xlink"
    viewBox="0 0 512 512"
    xmlSpace="preserve"
    fill="#ffffff"
    stroke="#ffffff"
    {...props}
  >
    <g id="SVGRepo_bgCarrier" strokeWidth={0} />
    <g
      id="SVGRepo_tracerCarrier"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <g id="SVGRepo_iconCarrier">
      <style type="text/css">{" .st0{fill:#ffffff;} "}</style>
      <g>
        <path
          className="st0"
          d="M194.651,414.476c16.23,0,61.349,0,61.349,0s45.111,0,61.35,0c16.222,0,23.587-23.603,14.198-40.285 c-7.072-12.572-18.659-26.826-37.516-34.921c-10.793,7.556-23.905,12-38.032,12c-14.143,0-27.238-4.444-38.032-12 c-18.864,8.095-30.444,22.349-37.523,34.921C171.064,390.873,178.421,414.476,194.651,414.476z"
        />
        <path
          className="st0"
          d="M256,335.476c27.714,0,50.167-22.444,50.167-50.159v-12.016c0-27.714-22.452-50.174-50.167-50.174 c-27.714,0-50.174,22.46-50.174,50.174v12.016C205.826,313.032,228.286,335.476,256,335.476z"
        />
        <path
          className="st0"
          d="M404.977,56.889h-75.834v16.254c0,31.365-25.524,56.889-56.889,56.889h-32.508 c-31.365,0-56.889-25.524-56.889-56.889V56.889h-75.833c-25.445,0-46.072,20.627-46.072,46.071v362.969 c0,25.444,20.627,46.071,46.072,46.071h297.952c25.444,0,46.071-20.627,46.071-46.071V102.96 C451.048,77.516,430.421,56.889,404.977,56.889z M402.286,463.238H109.714V150.349h292.572V463.238z"
        />
        <path
          className="st0"
          d="M239.746,113.778h32.508c22.406,0,40.635-18.23,40.635-40.635V40.635C312.889,18.23,294.659,0,272.254,0 h-32.508c-22.405,0-40.635,18.23-40.635,40.635v32.508C199.111,95.547,217.341,113.778,239.746,113.778z M231.619,40.635 c0-4.492,3.635-8.127,8.127-8.127h32.508c4.493,0,8.127,3.635,8.127,8.127v16.254c0,4.492-3.634,8.127-8.127,8.127h-32.508 c-4.492,0-8.127-3.635-8.127-8.127V40.635z"
        />
      </g>
    </g>
  </svg>
);

const Customer = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" {...props}>
    <path
      fill="#ffffff"
      d="M224 256A128 128 0 1 0 224 0a128 128 0 1 0 0 256zm-45.7 48C79.8 304 0 383.8 0 482.3C0 498.7 13.3 512 29.7 512H418.3c16.4 0 29.7-13.3 29.7-29.7C448 383.8 368.2 304 269.7 304H178.3z"
    />
  </svg>
);

const Avatar = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" {...props}>
    <path
      fill="#fff"
      d="M399 384.2c-22.1-38.4-63.6-64.2-111-64.2h-64c-47.4 0-88.9 25.8-111 64.2 35.2 39.2 86.2 63.8 143 63.8s107.8-24.7 143-63.8zM0 256a256 256 0 1 1 512 0 256 256 0 1 1-512 0zm256 16a72 72 0 1 0 0-144 72 72 0 1 0 0 144z"
    />
  </svg>
);

const Logout = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" {...props}>
    <path
      fill="#fff"
      d="m377.9 105.9 122.8 122.8c7.2 7.2 11.3 17.1 11.3 27.3s-4.1 20.1-11.3 27.3L377.9 406.1c-6.4 6.4-15 9.9-24 9.9-18.7 0-33.9-15.2-33.9-33.9V320H192c-17.7 0-32-14.3-32-32v-64c0-17.7 14.3-32 32-32h128v-62.1c0-18.7 15.2-33.9 33.9-33.9 9 0 17.6 3.6 24 9.9zM160 96H96c-17.7 0-32 14.3-32 32v256c0 17.7 14.3 32 32 32h64c17.7 0 32 14.3 32 32s-14.3 32-32 32H96c-53 0-96-43-96-96V128c0-53 43-96 96-96h64c17.7 0 32 14.3 32 32s-14.3 32-32 32z"
    />
  </svg>
);

const WarningIcon = (props) => (
  <svg
    viewBox="0 0 24 24"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    {...props}
  >
    <g id="SVGRepo_bgCarrier" strokeWidth={0} />
    <g
      id="SVGRepo_tracerCarrier"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <g id="SVGRepo_iconCarrier">
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M19.5 12C19.5 16.1421 16.1421 19.5 12 19.5C7.85786 19.5 4.5 16.1421 4.5 12C4.5 7.85786 7.85786 4.5 12 4.5C16.1421 4.5 19.5 7.85786 19.5 12ZM21 12C21 16.9706 16.9706 21 12 21C7.02944 21 3 16.9706 3 12C3 7.02944 7.02944 3 12 3C16.9706 3 21 7.02944 21 12ZM11.25 13.5V8.25H12.75V13.5H11.25ZM11.25 15.75V14.25H12.75V15.75H11.25Z"
        fill="#DC2626"
      />
    </g>
  </svg>
);

const SuccessIcon = (props) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    fill="#16A34A"
    stroke="#16A34A"
    viewBox="0 0 24 24"
    {...props}
  >
    <path
      fillRule="evenodd"
      d="M12 2c5.523 0 10 4.477 10 10s-4.477 10-10 10S2 17.523 2 12 6.477 2 12 2Zm0 2a8 8 0 1 0 0 16 8 8 0 0 0 0-16Zm3.293 4.293L10 13.586l-1.293-1.293a1 1 0 1 0-1.414 1.414l2 2a1 1 0 0 0 1.414 0l6-6a1 1 0 1 0-1.414-1.414Z"
    />
  </svg>
);

const FailureIcon = (props) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    xmlSpace="preserve"
    width={200}
    height={200}
    fill="#DC2626"
    viewBox="0 0 512 512"
    {...props}
  >
    <path d="M437.016 74.984c-99.979-99.979-262.075-99.979-362.033.002-99.978 99.978-99.978 262.073.004 362.031 99.954 99.978 262.05 99.978 362.029-.002 99.979-99.956 99.979-262.051 0-362.031zm-30.168 331.86c-83.318 83.318-218.396 83.318-301.691.004-83.318-83.299-83.318-218.377-.002-301.693 83.297-83.317 218.375-83.317 301.691 0s83.316 218.394.002 301.689z" />
    <path d="M361.592 150.408c-8.331-8.331-21.839-8.331-30.17 0l-75.425 75.425-75.425-75.425c-8.331-8.331-21.839-8.331-30.17 0s-8.331 21.839 0 30.17l75.425 75.425L150.43 331.4c-8.331 8.331-8.331 21.839 0 30.17 8.331 8.331 21.839 8.331 30.17 0l75.397-75.397 75.419 75.419c8.331 8.331 21.839 8.331 30.17 0 8.331-8.331 8.331-21.839 0-30.17l-75.419-75.419 75.425-75.425c8.331-8.331 8.331-21.838 0-30.17z" />
  </svg>
);

const Arrow = (props) => (
  <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 448 512" {...props}>
    <path d="M201.4 137.4c12.5-12.5 32.8-12.5 45.3 0l160 160c12.5 12.5 12.5 32.8 0 45.3s-32.8 12.5-45.3 0L224 205.3 86.6 342.6c-12.5 12.5-32.8 12.5-45.3 0s-12.5-32.8 0-45.3l160-160z" />
  </svg>
);
export const DashboardIcon = (props) => (
  <svg {...props}>
    <path fillRule="evenodd" d="M5.5 3.25...V15Z" />
  </svg>
);

export {
  Menu,
  Arrow,
  WarningIcon,
  SuccessIcon,
  FailureIcon,
  CheckFailure,
  Buses,
  Avatar,
  Logout,
  Route,
  Trip,
  BusStation,
  Ticket,
  Employee,
  Customer,
  
};
