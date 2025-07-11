/** @format */

import '../App.css';
import Banner from '../components/Banner';
import HomeBusLine from '../components/HomeBusLine';
import HomeFutaInfo from '../components/HomeFutaInfo';
import HomeNews from '../components/HomeNews';
import HomePromotion from '../components/HomePromotion';
import HomeSearch from '../components/HomeSearch';

function HomePage() {
	return (
		<div className="App px-2">
			<Banner />
			<HomeSearch />
			<HomePromotion />
			<HomeBusLine />
			<HomeFutaInfo />
			<HomeNews />
		</div>
	);
}

export default HomePage;
