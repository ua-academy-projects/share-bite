import "./App.css";
import { Routes, Route, NavLink } from "react-router-dom";
import { useEffect, useState } from "react";

type Box = {
  id: number;
  image: string;
  full_price: string;
  discount_price: string;
};

function Home() {
  return (
    <div className="p-8">
      <h1 className="text-4xl font-semibold text-white mb-2">
        Home Feed 🔥
      </h1>
      <p className="text-gray-400">
        Exceptional culinary creations from local ateliers.
      </p>
    </div>
  );
}

function Discover() {
  return <h1 className="text-white p-8">Discover 🌍</h1>;
}

function Boxes() {
  const [boxes, setBoxes] = useState<Box[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(
      "http://localhost:3999/business/nearby-boxes?lat=49.8397&lon=24.0297"
    )
      .then((res) => res.json())
      .then((data) => {
        setBoxes(data.items);
        setLoading(false);
      })
      .catch((err) => {
        console.error(err);
        setLoading(false);
      });
  }, []);

  if (loading) {
    return <div className="p-8 text-white">Loading...</div>;
  }

  if (boxes.length === 0) {
    return <div className="p-8 text-gray-400">No boxes found</div>;
  }

  return (
    <div className="p-8">
      <h1 className="text-4xl font-semibold text-white mb-2">
        Curated Rescues
      </h1>
      <p className="text-gray-400 mb-6">
        Exceptional culinary creations from local ateliers.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {boxes.map((box) => (
          <div
            key={box.id}
            className="bg-[#111] rounded-2xl overflow-hidden shadow-lg hover:scale-[1.02] transition"
          >
            <img
              src={
                box.image.includes("unsplash")
                  ? box.image.split("amazonaws.com/")[1]
                  : box.image}
              className="w-full h-40 object-cover"
            />

            <div className="p-4">
              <h3 className="text-white text-lg font-semibold">
                Surprise Box
              </h3>

              <div className="flex justify-between items-center mt-3">
                <span className="text-white font-semibold">
                  ${box.discount_price}
                </span>

                <button className="bg-green-500 text-black px-3 py-1 rounded-lg text-sm hover:bg-green-400 transition">
                  Reserve
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function App() {
  return (
    <div className="flex min-h-screen bg-[#0b0f0e]">
      <aside className="w-64 bg-[#163d32] p-6 flex flex-col justify-between">
        <div>
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 bg-black rounded-full flex items-center justify-center text-white">
              SB
            </div>
            <div>
              <h2 className="text-white font-semibold">Share Bite</h2>
              <p className="text-gray-300 text-sm">
                The Art of Dining
              </p>
            </div>
          </div>

          <button className="w-full bg-green-500 text-black py-2 rounded-full mb-8 hover:bg-green-400 transition">
            + Share a Bite
          </button>

          <nav className="flex flex-col gap-2">
            <NavLink
              to="/"
              end
              className={({ isActive }) =>
                `px-3 py-2 rounded-lg ${
                  isActive
                    ? "bg-[#2f5e50] text-white"
                    : "text-gray-200 hover:bg-[#2f5e50]"
                }`
              }
            >
              Home Feed
            </NavLink>

            <NavLink
              to="/boxes"
              className={({ isActive }) =>
                `px-3 py-2 rounded-lg ${
                  isActive
                    ? "bg-[#2f5e50] text-white"
                    : "text-gray-200 hover:bg-[#2f5e50]"
                }`
              }
            >
              Magic Boxes
            </NavLink>

            <NavLink
              to="/discover"
              className={({ isActive }) =>
                `px-3 py-2 rounded-lg ${
                  isActive
                    ? "bg-[#2f5e50] text-white"
                    : "text-gray-200 hover:bg-[#2f5e50]"
                }`
              }
            >
              Discover
            </NavLink>

            <span className="text-gray-300 px-3 py-2">
              Social Bites
            </span>
            <span className="text-gray-300 px-3 py-2">
              Settings
            </span>
          </nav>
        </div>

        <div className="text-gray-400 text-sm">
          <p>Support</p>
          <p>Privacy</p>
        </div>
      </aside>

      <main className="flex-1">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/boxes" element={<Boxes />} />
          <Route path="/discover" element={<Discover />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;