import Link from "next/link";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faList, faMagnifyingGlass, faPlus} from "@fortawesome/free-solid-svg-icons";
import styles from "../../public/navbar.module.css";

const Navbar = () => {
    return (
        <nav className={styles.Nav}>
            <Link href={"/"} className={styles.Link}>
                <FontAwesomeIcon icon={faList} />
            </Link>
            <Link href={"/find"} className={styles.Link}>
                <FontAwesomeIcon icon={faMagnifyingGlass}/>
            </Link>
            <Link href={"/calculate"} className={styles.Link}>
                <FontAwesomeIcon icon={faPlus}/>
            </Link>
        </nav>
    )
}

export default Navbar